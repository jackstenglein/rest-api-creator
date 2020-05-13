import React from 'react';
import ProjectDetails from './ProjectDetails.js';
import { getDownloadLink } from '../api/api.js';
import { saveAs } from 'file-saver';

class ProjectDetailsContainer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};
    this.downloadCode = this.downloadCode.bind(this);
  }

  async downloadCode() {
    console.log("Making request to get download URL");
    this.setState({downloadError: undefined, downloadStatus: "Generating Code"});
    var response = await getDownloadLink(this.props.project.id);
    console.log("Got download code response: ", response);
    if (response.error !== undefined) {
      this.setState({downloadError: response.error, downloadStatus: undefined});
    } else {
      this.setState({downloadStatus: undefined});
      saveAs(response.url, "defaultProject.zip");
    }
  }

  render() {
    return (
      <ProjectDetails 
        project={this.props.project} 
        downloadStatus={this.state.downloadStatus}
        downloadError={this.state.downloadError}
        download={this.downloadCode}
      />
    );
  }
}

export default ProjectDetailsContainer;
