import React from 'react';
import { connect } from 'react-redux';
import { STATUS_PENDING, STATUS_NONE } from '../../redux/modules/network';
import ObjectDetailView from './ObjectDetailView';
import { deleteObject } from '../../api/api';
import { deleteObjectSuccess } from '../../redux/modules/projects'; 
import produce from 'immer';
import { Redirect } from 'react-router-dom';

class ObjectDetailContainer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};
    this.delete = this.delete.bind(this);
  }

  async delete() {
    const projectId = this.props.project.id;
    const objectId = this.props.match.params.objectId;
    const response = await deleteObject(projectId, objectId);

    var apiError = response.error;
    var deleted = response.error === undefined;
    if (deleted) {
      this.props.onDelete(projectId, objectId);
    }

    const nextState = produce(this.state, draftState => {
      draftState.apiError = apiError;
      draftState.deleted = deleted;
    });
    this.setState(nextState);
  }

  render() {
    if (this.props.project.network.status === STATUS_NONE || this.props.project.network.status === STATUS_PENDING) {
      return null;
    }

    if (this.state.deleted) {
      return <Redirect to="/app/objects"/>
    }

    const objectId = this.props.match.params.objectId;
    const object = this.props.project.objects[objectId];
    const error = this.state.apiError || this.props.project.network.error;

    return (
      <ObjectDetailView object={object} error={error} delete={this.delete}/>
    )
  }
}

const mapDispatchToProps = dispatch => {
  return {
    onDelete: (projectId, objectId) => {
      dispatch(deleteObjectSuccess(projectId, objectId))
    }
  }
}

export default connect(null, mapDispatchToProps)(ObjectDetailContainer);
