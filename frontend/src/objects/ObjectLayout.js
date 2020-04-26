import React from 'react';
import { Route, Switch } from 'react-router-dom';
import ObjectEditor from './ObjectEditor.js';
import ObjectList from './ObjectList.js';

const ObjectLayout = props => (
  <div className="object-layout">
    <Switch>
      <Route path={props.match.path} exact component={ObjectList} />
      <Route path={`${props.match.path}:objectId`} render={(props) => 
        <ObjectEditor {...props} values={EDITOR_VALUES} errors={EDITOR_ERRORS} path={EDITOR_PATH}/>} 
      />
    </Switch>
  </div>
)

const EDITOR_PATH = [
  {
    name: "Default Project",
    link: "/app",
    active: false
  },
  {
    name: "Objects",
    link: "/app/objects",
    active: false
  }, 
  {
    name: "User",
    link: "/app/objects/5",
    active: false
  }, 
  {
    name: "Edit",
    link: "/app/objects/5/edit",
    active: true
  }
]

const EDITOR_VALUES = {
  name: "User",
  description: "User represents the structure of an end user in the database.",
  attributes: [
    {
      name: "AttributeOne",
      type: "Integer",
      defaultValue: "5",
    },
    {
      name: "",
      type: "",
      defaultValue: "",
    }
  ]
}

const EDITOR_ERRORS = {
  name: "Object with this name already exists.",
  attributes: [
    {},
    {
      name: "Name is required",
      type: "Type is required"
    }
  ]
}

export default ObjectLayout;
