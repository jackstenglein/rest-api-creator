import React from 'react';
import { STATUS_PENDING, STATUS_NONE } from '../../redux/modules/network';
import ObjectDetailView from './ObjectDetailView';


class ObjectDetailContainer extends React.Component {

    constructor(props) {
        super(props);
        this.delete = this.delete.bind(this);
    }

    delete() {
        console.log('Delete object');
    }

    render() {
        if (this.props.project.network.status === STATUS_NONE || this.props.project.network.status === STATUS_PENDING) {
            return null;
        }

        const objectId = this.props.match.params.objectId;
        const object = this.props.project.objects[objectId];
        const error = this.props.project.network.error;

        return (
            <ObjectDetailView object={object} error={error} delete={this.delete}/>
        )
    }
}

export default ObjectDetailContainer;
