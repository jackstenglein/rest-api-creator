import React from 'react';
import { Route } from 'react-router-dom';
import { connect } from 'react-redux';
import { STATUS_NONE, STATUS_FAILURE, none } from '../redux/modules/network.js';
import { fetchProjectRequest, fetchProjectResponse } from '../redux/modules/projects';
import { getProject } from '../api/api';


class ProjectRoute extends React.Component {

  async componentDidMount() {
    console.log("Project Route component did mount:  ", this.props.project);
    const network = this.props.project.network;
    if (network.status === STATUS_NONE || network.status === STATUS_FAILURE) {
      console.log("Getting the project");
      const projectId = "defaultProject"; // TODO: dynamically get this from URL
      this.props.getProjectRequest(projectId);
      const response = await getProject(projectId);
      this.props.getProjectResponse(projectId, response);
    }
  }

  render() {
    const { component: Component, project, ...rest } = this.props;
    return (
      <Route {...rest} render={props => (<Component {...props} project={project}/>)} />
    )
  }
}

const mapStateToProps = state => {
  const projectId = "defaultProject"; // TODO: dynamically get this from URL
  const project = state.projects[projectId]; 
  if (project !== undefined) {
    return {project: project};
  }

  return {
    project: { 
      network: none()
    }
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

export default connect(mapStateToProps, mapDispatchToProps)(ProjectRoute);
