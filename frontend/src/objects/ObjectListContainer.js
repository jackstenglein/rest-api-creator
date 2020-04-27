import { connect } from 'react-redux';
import ObjectListView from './ObjectListView';

const mapStateToProps = state => {
  return {
    objects: state.projects.defaultProject.objects
  };
}

export default connect(mapStateToProps)(ObjectListView);
