import React, {Component} from 'react';

class FilePicker extends Component {

  uploadFile = e => {
    if (!e.currentTarget.files.length) {
      return;
    }

    const file = e.currentTarget.files[0];

    if (!file.type.match(/^image\//g)) {
      alert('Image files only');
      return;
    }

    this.props.fileSelected(file);
  };

  render() {
    return (
        <label htmlFor="upload" className="btn btn-primary">
          <input type="file" style={{display: 'none'}}
                 onChange={this.uploadFile} id="upload"/>
          Select image to shrinkflate
        </label>
    );
  }
}

export default FilePicker;
