import React from 'react';
import { Route, Switch } from 'react-router-dom';
import ObjectEditor from './ObjectEditor.js';
import ObjectListView from './ObjectListView.js';

const ObjectLayout = props => (
  <div className="object-layout">
    <Switch>
      <Route path={props.match.path} exact render={(props) => <ObjectListView {...props} objects={TEST_OBJECTS} />} />
      <Route path={`${props.match.path}create`} render={props =>
        <ObjectEditor {...props} values={EMPTY_OBJECT} errors={EMPTY_OBJECT} />} />
      <Route path={`${props.match.path}:objectId/edit`} render={(props) => 
        <ObjectEditor {...props} values={EDITOR_VALUES} errors={EDITOR_ERRORS} path={EDITOR_PATH}/>} 
      />
      
    </Switch>
  </div>
)

const EMPTY_OBJECT = {
  name: "",
  description: "",
  attributes: []
}

const TEST_OBJECTS = [
  {
    name: "User",
    description: "User represents the structure of an end user in the database.",
    attributes: [
      {
        name: "AttributeOne",
        type: "Integer",
        required: true,
        description: "AttributeOne is an integer attribute that stores the first attribute of the user."
      },
      {
        name: "AttributeTwo",
        type: "Text",
        required: false,
        description: "AttributeTwo is a text attribute that stores the second attribute of the user."
      }
    ]
  },
  {
    name: "Dog",
    description: "Dog represents the structure of a dog object in the database.",
    attributes: [
      {
        name: "AttributeOne",
        type: "Text",
        required: true,
        description: "AttributeOne is a text attribute that stores the first attribute of the dog. This attribute is required."
      },
      {
        name: "AttributeTwo",
        type: "Text",
        required: true,
        description: "AttributeTwo is a text attribute that stores the second attribute of the dog."
      }
    ]
  }
]

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
