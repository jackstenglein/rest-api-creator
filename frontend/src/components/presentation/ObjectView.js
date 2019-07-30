import React, { Component } from 'react'
import Redirect from 'react-router-dom/Redirect';
import Container from 'react-bootstrap/Container';
import ListGroup from 'react-bootstrap/ListGroup';
import Breadcrumb from 'react-bootstrap/Breadcrumb';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Spinner from 'react-bootstrap/Spinner';
import Table from 'react-bootstrap/Table';
// import {
//     CLICK_CREATE,
//     FETCH_OBJECTS_REQUEST
// } from '../../redux/actions/objects/objectListActions.js';
const ID_DESCRIPTION = "The id attribute uniquely identifies database records of this"
    + " object type. There can only be one record with a given id. Each time an"
    + " object of this type is created, it is automatically given a new id.";

const testObject = {
    details: {
        'name': 'Test Object Name',
        'description': 'This is a test object description. It is hopefully long enough to span across two lines.'
    },
    attributes: [
        {
            'name': 'name',
            'type': 'Text',
            'description': 'The name of the user.'
        }
    ]
};

class ObjectView extends Component {
    constructor(props) {
        super(props);
        this.handleRefreshClick = this.handleRefreshClick.bind(this);
    }

    componentDidMount() {
        console.log("ObjectView did mount");
        // const { dispatch, selectedSubreddit } = this.props
        // dispatch(fetchPostsIfNeeded(selectedSubreddit))
        // this.props.callbacks.reset();
        // this.props.callbacks.fetchDetailsIfNeeded();
    }
    //
    // componentDidUpdate(prevProps) {
    //     console.log("Component did update. Do we need this?");
    //     this.props.callbacks.reset();
    // }

  // handleChange(nextSubreddit) {
  //   this.props.dispatch(selectSubreddit(nextSubreddit))
  //   this.props.dispatch(fetchPostsIfNeeded(nextSubreddit))
  // }

    handleRefreshClick(event) {
        event.preventDefault()

        // const { dispatch, selectedSubreddit } = this.props
        // dispatch(invalidateSubreddit(selectedSubreddit))
        // dispatch(fetchPostsIfNeeded(selectedSubreddit))
    }

    getIdAttributeJsx() {
        return (
            <tr>
                <td>id</td>
                <td>Integer</td>
                <td>{ID_DESCRIPTION}</td>
            </tr>
        );
    }

    getBreadcrumbJsx() {
        return (
            <Breadcrumb bsPrefix="breadcrumb crud">
                <Breadcrumb.Item href="#">Projects</Breadcrumb.Item>
                <Breadcrumb.Item href="#">{this.props.projectName}</Breadcrumb.Item>
                <Breadcrumb.Item href="#">Objects</Breadcrumb.Item>
            </Breadcrumb>
        );
    }

    getTopBar() {
        let editButton;
        let deleteButton;
        if (false) {
            editButton = (<Button variant="secondary" className="mr-2" disabled />);
            deleteButton = (
                <Button variant="danger" disabled>
                    <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true"/>
                </Button>
            );
        } else {
            editButton = (<Button variant="secondary" className="mr-2">Edit</Button>);
            deleteButton = (<Button variant="danger">Delete</Button>);
        }

        return (
            <Row className="align-items-center justify-content-between object-editor-toolbar">
                <Col xs="auto">
                    { this.getBreadcrumbJsx() }
                </Col>
                <Col xs="auto">
                    <ButtonToolbar>
                        { editButton }
                        { deleteButton }
                    </ButtonToolbar>
                </Col>
            </Row>
        );
    }

    render() {
        console.log("Object view props: ", this.props);

        // const { selectedSubreddit, posts, isFetching, lastUpdated } = this.props
        return (
            <Container className="object-editor-viewport">
                { this.getTopBar() }
                <h5>{testObject.details.name}</h5>
                <p>{testObject.details.description}</p>
                <Table bordered striped>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Data Type</th>
                            <th>Description</th>
                        </tr>
                    </thead>
                    <tbody>
                        { this.getIdAttributeJsx() }
                        {
                            testObject.attributes.map(function(attribute) {
                                return (
                                    <tr>
                                        <td>{attribute.name}</td>
                                        <td>{attribute.type}</td>
                                        <td>{attribute.description}</td>
                                    </tr>
                                )
                            })
                        }
                    </tbody>
                </Table>
            </Container>
        );
    }
}

export default ObjectView;
