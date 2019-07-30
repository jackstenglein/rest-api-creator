import React from 'react'
import PropTypes from 'prop-types'
import { Provider } from 'react-redux'
import { BrowserRouter as Router, Route } from 'react-router-dom'
import App from './App'
import ObjectEditorContainer from './components/container/ObjectEditorContainer';
import ObjectListViewContainer from './components/container/ObjectListViewContainer';
import ObjectViewContainer from './components/container/ObjectViewContainer';


const Root = ({ store }) => (
  <Provider store={store}>
  <link
rel="stylesheet"
href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
crossorigin="anonymous"
/>
    <Router>
      <Route exact path="/" component={App} />
      <Route exact path="/objects" component={ObjectListViewContainer} />
      <Route path="/objects/:id" component={ObjectViewContainer} />
      <Route path="/createObject" component={ObjectEditorContainer} />
    </Router>
  </Provider>
)

Root.propTypes = {
  store: PropTypes.object.isRequired
}

export default Root;
