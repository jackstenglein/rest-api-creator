import React from 'react';
import Container from 'react-bootstrap/Container';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import ProjectMenu from './ProjectMenu';
import ObjectListViewContainer from '../container/ObjectListViewContainer';
import ObjectDetailsContainer from '../container/ObjectDetailsContainer';
import ObjectEditorContainer from '../container/ObjectEditorContainer';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

function ObjectView(props) {
    return (
            <Row>
                <Col className='project-menu' xl={2}><ProjectMenu selectedId='objects'/></Col>
                <Col>
                    <Router>
                        <Switch>
                            <Route exact path="/objects" component={ObjectListViewContainer} />
                            <Route exact path="/objects/create" component={ObjectEditorContainer} />
                            <Route path="/objects/:id" component={ObjectDetailsContainer} />
                        </Switch>
                    </Router>
                </Col>
            </Row>
    );
}

export default ObjectView;
