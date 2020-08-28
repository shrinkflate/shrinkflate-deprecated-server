import React, {Component} from 'react';
import FileSize from '../Filesize';
import FilePicker from '../FilePicker';

class Source extends Component {

  state = {
    size: 0,
    image: null,
  };

  fileSelected = file => {
    this.setState({size: file.size});
    const reader = new FileReader();
    reader.onload = e => {
      this.setState({image: e.target.result});
    };
    reader.readAsDataURL(file);
    this.props.fileSelected(file);
  };

  render() {
    return (
        <div className="row">
          <div className="col-md-12 d-flex justify-content-between">
            <FileSize title="Original image" size={this.state.size}/>
            <FilePicker fileSelected={this.fileSelected}/>
          </div>
          <div className="col-md-12">
            {this.state.image ?
                <img src={this.state.image} className="img-fluid" alt=""/> :
                null}
          </div>
        </div>
    );
  }
}

export default Source;
