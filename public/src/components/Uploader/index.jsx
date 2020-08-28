import React, {Component} from 'react';
import axios from 'axios';

class Uploader extends Component {

  state = {
    compressor: 'libvips',
    quality: 85,
    progressive: false,
    uploading: false,
  };

  componentDidMount() {
    this.setState({
      compressor: this.props.opts.compressor || 'libvips',
      quality: this.props.opts.quality || 85,
      progressive: this.props.opts.progressive || false,
    });
  }

  setCompressor = e => {
    this.setState({compressor: e.currentTarget.value});
  };

  setQuality = e => {
    this.setState({quality: e.currentTarget.value});
  };

  setProgressive = e => {
    this.setState({progressive: e.currentTarget.checked});
  };

  startUpload = () => {
    this.setState({uploading: true});
    this.props.setStatus('uploading');

    const form = new FormData();
    form.append('image', this.props.file);
    form.append('quality', this.state.quality);
    form.append('progressive', this.state.progressive);
    form.append('compressor', this.state.compressor);

    const prefix = window.location.href.includes('localhost') ?
        'http://localhost:4000' :
        '';

    axios.post(`${prefix}/compress`, form).then(response => {
      this.props.uploaded({
        fileId: response.data,
        opts: {
          compressor: this.state.compressor,
          quality: this.state.quality,
          progressive: this.state.progressive,
        },
      });
    }).catch(error => {
      this.setState({uploading: false});
      alert('Error occurred. check console log.');
      console.log(error);
    });
  };

  render() {
    return (
        <div className={`col-md-12 ${this.props.className}`}>
          <div className="form-group row">
            <label htmlFor="compressor"
                   className="col-sm-4 col-form-label">Compressor</label>
            <div className="col-sm-6">
              <select value={this.state.compressor}
                      onChange={this.setCompressor}
                      className="form-control"
                      id="compressor"
              >
                <option value="libvips">libvips(vips)</option>
                <option value="lilliput">lilliput(OpenCV)</option>
              </select>
            </div>
          </div>
          {this.state.compressor === 'lilliput'
              ? [
                <div className="form-group row" key="quality">
                  <label
                      htmlFor="quality"
                      className="col-sm-4 col-form-label"
                  >Quality</label>
                  <div className="col-sm-6">
                    <input
                        type="number"
                        min={30}
                        max={100}
                        value={this.state.quality}
                        onChange={this.setQuality}
                        className="form-control"
                        id="quality"/>
                  </div>
                </div>,

                <div className="form-group row" key="progressive">
                  <div className="col-sm-6 offset-4 form-check">
                    <input
                        type="checkbox"
                        checked={this.state.progressive}
                        onChange={this.setProgressive}
                        id="progressive"
                        className="form-check-input"
                    />

                    <label
                        htmlFor="progressive"
                        className="form-check-label"
                    >Progressive</label>
                  </div>
                </div>,
              ]
              : null}

          <button className="btn btn-primary" disabled={this.state.uploading}
                  onClick={this.startUpload}>Compress
          </button>
        </div>
    );
  }
}

export default Uploader;
