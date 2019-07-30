import React, { Component } from 'react'
import { connect } from 'react-redux'
import ObjectListView from '../presentation/ObjectListView';
import {
    clickCreate,
    clickObject,
    fetchObjects,
    fetchObjectsIfNeeded,
    resetListView
} from '../../redux/actions/objects/objectListActions.js';



function mapStateToProps(state) {
    return Object.assign({}, state.objects.list, {projectName: 'Test Project'});
}

const mapDispatchToProps = dispatch => {
    return {
        callbacks: {
            reset: () => dispatch(resetListView()),
            fetchObjectsIfNeeded: () => dispatch(fetchObjectsIfNeeded()),
            refresh: () => dispatch(fetchObjects()),
            onClickCreate: () => dispatch(clickCreate()),
            clickObject: (id) => dispatch(clickObject(id))
        }
    };
}

const ObjectListViewContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectListView);

export default ObjectListViewContainer;
