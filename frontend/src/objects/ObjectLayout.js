import React from 'react';
import { Route, Switch } from 'react-router-dom';
import ObjectEditorContainer from './editor/ObjectEditorContainer.js';
import ObjectListContainer from './ObjectListContainer.js';

const ObjectLayout = props => (
  <div className="object-layout">
    <Switch>
      <Route path={props.match.path} exact component={ObjectListContainer} />
      <Route path={`${props.match.path}create`} component={ObjectEditorContainer} />
      <Route path={`${props.match.path}:objectId/edit`} component={ObjectEditorContainer} />
    </Switch>
  </div>
)

export default ObjectLayout;
