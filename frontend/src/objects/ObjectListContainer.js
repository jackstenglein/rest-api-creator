import React from 'react';
import { connect } from 'react-redux';
import ObjectListView from './ObjectListView';
import * as network from '../redux/modules/network.js';
import { getProject } from '../api/api.js';
import { fetchProjectRequest, fetchProjectResponse } from '../redux/modules/projects';

class ObjectListContainer extends React.Component {

  constructor(props) {
    super(props);
    this.getProject = this.getProject.bind(this);
  }

  async componentDidMount() {
    console.log("ObjectListContainer props: ", this.props);
    if (this.props.network.status === network.STATUS_NONE || this.props.network.status === network.STATUS_FAILURE) {
      console.log("Making get project request");
      this.getProject();
    }
  }

  async getProject() {
    const projectId = "defaultProject"; // TODO: dynamically get this from URL
    this.props.getProjectRequest(projectId);
    const response = await getProject(projectId);
    console.log("Get project response: ", response);
    this.props.getProjectResponse(projectId, response);
  }

  render() {
    return <ObjectListView {...this.props} onRefresh={this.getProject}/>
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
    objects: project.objects
  };
}

const mapDispatchToProps = dispatch => {
  return {
    getProjectRequest: (id) => {
      dispatch(fetchProjectRequest(id))
    },
    getProjectResponse: (id, response) => {
      dispatch(fetchProjectResponse(id, response))
    }
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ObjectListContainer);
