import React, {Component} from 'react';

class FileSize extends Component {

  get size() {
    if (this.props.size) {
      return this.fileSize(this.props.size);
    }
    return '0B';
  }

  render() {
    return (
        <h3>{this.props.title}: {this.size}</h3>
    );
  }

  fileSize(size) {
    const i = Math.floor(Math.log(size) / Math.log(1024));
    return (size / Math.pow(1024, i)).toFixed(2) * 1 + ' ' +
        ['B', 'kB', 'MB', 'GB', 'TB'][i];
  }
}

export default FileSize;