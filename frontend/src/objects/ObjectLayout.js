import React from 'react';
import { Switch } from 'react-router-dom';
import ObjectEditorContainer from './editor/ObjectEditorContainer.js';
import ObjectListView from './ObjectListView.js';
import ProjectRoute from '../navigation/ProjectRoute.js';

const ObjectLayout = props => (
  <Switch>
    <ProjectRoute path={props.match.path} exact component={ObjectListView} />
    <ProjectRoute path={`${props.match.path}create`} component={ObjectEditorContainer} />
    <ProjectRoute path={`${props.match.path}:objectId/edit`} component={ObjectEditorContainer} />
  </Switch>
)

export default ObjectLayout;
