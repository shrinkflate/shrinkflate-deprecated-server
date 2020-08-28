import React, {Component} from 'react';

class PageHeader extends Component {
  render() {
    return (
        <div className="col-md-12">
          <h2 className="text-center">
            <span>Shrinkflate</span>
            {
              this.props.status
                  ? <small style={{fontSize: 'medium'}}> {this.props.status} </small>
                  : null
            }
          </h2>
        </div>
    );
  }
}

export default PageHeader;
