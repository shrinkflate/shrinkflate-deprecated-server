import React, {Component} from 'react';
import PageHeader from '../PageHeader';
import Source from '../Source';
import Destination from '../Destination';

class Container extends Component {

  state = {
    status: null,
    file: null,
  };

  setFileId = fileId => {
    this.setState({fileId});
  };

  setStatus = status => {
    this.setState({status});
  };

  fileSelected = file => {
    this.setState({file});
  };

  render() {
    return (
        <div className="container-fluid">
          <div className="row">
            <PageHeader status={this.state.status}/>
          </div>
          <div className="row">
            <div className="col-md-6">
              <Source
                  setFileId={this.setFileId}
                  fileSelected={this.fileSelected}
              />
            </div>
            <div className="col-md-6">
              <Destination
                  setStatus={this.setStatus}
                  file={this.state.file}
              />
            </div>
          </div>
        </div>
    );
  }
}

export default Container;