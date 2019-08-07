import React from 'react'
import PropTypes from 'prop-types'
import { Provider } from 'react-redux'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import App from './App'
import ObjectEditorContainer from './components/container/ObjectEditorContainer';
import ObjectListViewContainer from './components/container/ObjectListViewContainer';
import ObjectDetailsContainer from './components/container/ObjectDetailsContainer';
import ObjectView from './components/presentation/ObjectView';


const Root = ({ store }) => (
  <Provider store={store}>
  <link
rel="stylesheet"
href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
crossorigin="anonymous"
/>
    <Router>
        <Switch>
            <Route exact path="/" component={App} />
            <Route path="/objects" component={ObjectView} />
            {/* <Route exact path="/objects" component={ObjectListViewContainer} />
            <Route path="/objects/:id" component={ObjectDetailsContainer} />
            <Route path="/createObject" component={ObjectEditorContainer} /> */}
        </Switch>
    </Router>
  </Provider>
)

Root.propTypes = {
  store: PropTypes.object.isRequired
}

export default Root;
