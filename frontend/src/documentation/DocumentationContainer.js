import {connect} from 'react-redux';
import Documentation from './Documentation';

const TEST_PROJECT = {
  id: "defaultProject",
  name: "Default Project",
  description: "This is the default project. Future versions of this website will allow you to create multiple projects. " + 
               "If you had created your own project, a description of it would go here. For now, this site allows you to " +
               "define database objects and auto-generate the code to support a REST API based on those objects. Future " +
               "versions will also allow you to automatically deploy the API to AWS and other cloud services. This page " +
               "will allow you to see the status of those deployed APIs.",
  objects: {
  }
}

const mapStateToProps = state => {
  return {
    project: TEST_PROJECT//state.projects.defaultProject
  };
}

const DocumentationContainer = connect(mapStateToProps)(Documentation);
export default DocumentationContainer;
