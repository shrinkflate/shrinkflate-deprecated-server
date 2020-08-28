import React, {Component} from 'react';
import axios from 'axios';
import FileSize from '../Filesize';
import Uploader from '../Uploader';

class Destination extends Component {

  state = {
    file: null,
    size: 0,
    url: null,
    fileId: null,
    opts: {},
  };

  reconfigure = () => {
    this.setState({file: null, fileId: null, url: null, size: 0});
  };

  uploaded = e => {
    this.setState({file: this.props.file, fileId: e.fileId, opts: e.opts});
    this.props.setStatus('Uploaded');
    this.lookForFinish();
  };

  lookForFinish = () => {
    this.props.setStatus('waiting for finish...');

    const prefix = window.location.href.includes('localhost') ?
        'http://localhost:4000' :
        '';
    const url = `${prefix}/download/${this.state.fileId}`;
    axios.get(url + '?q=' + (new Date()).getTime()).then(response => {
      const size = response.headers['content-length'];
      this.setState({url, size});
      this.props.setStatus('done');
    }).catch(error => {
      if (error.response && error.response.status === 404) {
        setTimeout(this.lookForFinish, 2000);
      } else {
        alert('Error occurred. check console');
        console.log(error);
        this.props.setStatus('error occurred');
      }
    });
  };

  render() {
    return (
        <div className="row">
          {this.state.url && this.props.file === this.state.file
              ? <div className="col-md-12 d-flex justify-content-between">
                <FileSize title="Compressed image" size={this.state.size}/>
                <button className="btn btn-primary"
                        onClick={this.reconfigure}>Reconfigure
                </button>
              </div>
              : null}
          {this.state.url && this.props.file === this.state.file
              ? <div className="col-md-12">
                <img src={this.state.url} className="img-fluid" alt=""/>
              </div>
              : null}
          {this.props.file !== this.state.file
              ? <Uploader className="mt-5" uploaded={this.uploaded}
                          setStatus={this.props.setStatus}
                          file={this.props.file}
                          opts={this.state.opts}
              />
              : null}
        </div>
    );
  }
}

export default Destination;
