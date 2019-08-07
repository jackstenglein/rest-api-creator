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
import {
    CLICK_CREATE,
    FETCH_OBJECTS_REQUEST
} from '../../redux/actions/objects/objectListActions.js';

class ObjectListView extends Component {
    constructor(props) {
        super(props);
        this.handleRefreshClick = this.handleRefreshClick.bind(this);
    }

    componentDidMount() {
        console.log("Component did mount");
        // const { dispatch, selectedSubreddit } = this.props
        // dispatch(fetchPostsIfNeeded(selectedSubreddit))
        this.props.callbacks.reset();
        this.props.callbacks.fetchObjectsIfNeeded();
    }

    componentDidUpdate(prevProps) {
        console.log("Component did update. Do we need this?");
        this.props.callbacks.reset();
    }

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
        let refreshButton;
        let createButton;
        if (this.props.status === FETCH_OBJECTS_REQUEST) {
            refreshButton = (
                <Button variant="secondary" className="mr-2" disabled>
                    <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true"/>
                </Button>
            );
            createButton = (<Button variant="primary" disabled>Create</Button>);
        } else {
            refreshButton = (
                <Button variant="secondary" className="mr-2" onClick={this.props.callbacks.refresh}>Refresh</Button>
            );
            createButton = (<Button variant="primary" onClick={this.props.callbacks.onClickCreate}>Create</Button>);
        }

        return (
            <Row className="align-items-center justify-content-between object-editor-toolbar">
                <Col xs="auto">
                    { this.getBreadcrumbJsx() }
                </Col>
                <Col xs="auto">
                    <ButtonToolbar>
                        { refreshButton }
                        { createButton }
                    </ButtonToolbar>
                </Col>
            </Row>
        );
    }

    render() {
        console.log("Object list view props: ", this.props);
        if (this.props.create) {
            return (<Redirect to="/objects/create" />);
        }

        if (this.props.selectedObject !== -1) {
            return (<Redirect to={"/objects/" + this.props.selectedObject} />);
        }

        // const { selectedSubreddit, posts, isFetching, lastUpdated } = this.props
        return (
            <Container className="object-editor-viewport">
                { this.getTopBar() }
                <ListGroup variant="flush">
                    {
                        this.props.items.map(function(item) {
                            return (
                                <ListGroup.Item
                                    action
                                    key={item.id}
                                    onClick={() => this.props.callbacks.clickObject(item.id)}
                                >
                                    {item.name}
                                </ListGroup.Item>
                            );
                        }, this)
                    }
                </ListGroup>
            </Container>
        );
    }
}

export default ObjectListView;
