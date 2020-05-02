import React from 'react';
import {connect} from 'react-redux';
import Documentation from './Documentation';
import * as network from '../redux/modules/network.js';
import { getProject } from '../api/api';
import { fetchProjectRequest, fetchProjectResponse } from '../redux/modules/projects';


class DocumentationContainer extends React.Component {
  async componentDidMount() {
    if (this.props.network.status === network.STATUS_NONE || this.props.network.status === network.STATUS_FAILURE) {
      const projectId = "defaultProject"; // TODO: dynamically get this from URL
      this.props.getProjectRequest(projectId);
      const response = await getProject(projectId);
      this.props.getProjectResponse(projectId, response);
    }
  }

  render() {
    return (
      <Documentation project={this.props.project}/>
    )
  }
}

const mapStateToProps = state => {
  const project = state.projects["defaultProject"];
  if (project === undefined) {
    return {
      network: network.none()
    };
  }
  return {
    network: project.network,
    project: project
  };
}

const mapDispatchToProps = dispatch => {
  return {
    getProjectRequest: id => {
      dispatch(fetchProjectRequest(id))
    },
    getProjectResponse: (id, response) => {
      dispatch(fetchProjectResponse(id, response))
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(DocumentationContainer);
