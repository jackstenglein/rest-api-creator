import {connect} from 'react-redux';
import Documentation from './Documentation';

const mapStateToProps = state => {
  return {
    project: state.projects.defaultProject
  };
}

const DocumentationContainer = connect(mapStateToProps)(Documentation);
export default DocumentationContainer;
