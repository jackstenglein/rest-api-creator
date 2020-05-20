import React from 'react';
import { Switch } from 'react-router-dom';
import ObjectEditorContainer from './editor/ObjectEditorContainer.js';
import ObjectListView from './ObjectListView.js';
import ProjectRoute from '../navigation/ProjectRoute.js';
import ObjectDetailContainer from './detail/ObjectDetailContainer.js';

const ObjectLayout = props => (
  <Switch>
    <ProjectRoute path={props.match.path} exact component={ObjectListView} />
    <ProjectRoute path={`${props.match.path}:objectId/edit`} component={ObjectEditorContainer} />
    <ProjectRoute path={`${props.match.path}:objectId`} component={ObjectDetailContainer} />
    <ProjectRoute path={`${props.match.path}create`} component={ObjectEditorContainer} />
  </Switch>
)

export default ObjectLayout;
