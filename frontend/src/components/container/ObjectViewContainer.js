import { connect } from 'react-redux'
import ObjectView from '../presentation/ObjectView';
// import {
//     clickCreate,
//     clickObject,
//     fetchObjects,
//     fetchObjectsIfNeeded,
//     resetListView
// } from '../../redux/actions/objects/objectListActions.js';



function mapStateToProps(state, ownProps) {
    console.log("Own props: ", ownProps);
    return Object.assign({}, state.objects.list.items[31], {projectName: 'Test Project'});
}

const mapDispatchToProps = dispatch => {
    return {
        // callbacks: {
        //     reset: () => dispatch(resetListView()),
        //     fetchObjectsIfNeeded: () => dispatch(fetchObjectsIfNeeded()),
        //     refresh: () => dispatch(fetchObjects()),
        //     onClickCreate: () => dispatch(clickCreate()),
        //     clickObject: (id) => dispatch(clickObject(id))
        // }
    };
}

const ObjectViewContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(ObjectView);

export default ObjectViewContainer;
