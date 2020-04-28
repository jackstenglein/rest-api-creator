import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './navigation/App';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import store from './redux/store';
import { putObjectSuccess, putProjectSuccess } from './redux/modules/projects';

const TEST_PROJECT = {
  id: "defaultProject",
  name: "Default Project",
  description: "This is the default project. Future versions of this website will allow you to create multiple projects. " + 
               "If you had created your own project, a description of it would go here. For now, this site allows you to " +
               "define database objects and auto-generate the code to support a REST API based on those objects. Future " +
               "versions will also allow you to automatically deploy the API to AWS and other cloud services. This page " +
               "will allow you to see the status of those deployed APIs.",
  objects: {
    "user": {
      id: "user",
      name: "User",
      description: "User represents the structure of an end user in the database.",
      attributes: [
        {
          name: "AttributeOne",
          type: "Integer",
          required: true,
          description: "AttributeOne is an integer attribute that stores the first attribute of the user.",
          defaultValue: ""
        },
        {
          name: "AttributeTwo",
          type: "Text",
          required: false,
          description: "AttributeTwo is a text attribute that stores the second attribute of the user.",
          defaultValue: ""
        }
      ]
    },
    "dog": {
      id: "dog",
      name: "Dog",
      description: "Dog represents the structure of a dog object in the database.",
      attributes: [
        {
          name: "AttributeOne",
          type: "Text",
          required: true,
          description: "AttributeOne is a text attribute that stores the first attribute of the dog. This attribute is required.",
          defaultValue: "Asdf"
        },
        {
          name: "AttributeTwo",
          type: "Text",
          required: true,
          description: "AttributeTwo is a text attribute that stores the second attribute of the dog.",
          defaultValue: ""
        }
      ]
    }
  }
}

store.dispatch(putProjectSuccess(TEST_PROJECT))
store.dispatch(putObjectSuccess("defaultProject", {id: "objectname", name: "ObjectName", description: "Test object", attributes: []}))

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
