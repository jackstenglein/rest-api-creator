import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './navigation/App';
import { BrowserRouter } from 'react-router-dom';

// import { createStore, applyMiddleware } from 'redux';
// import { crudCreatorApp } from './redux/reducers/mainReducer.js';
// import thunkMiddleware from 'redux-thunk';
// import Root from './Root'

// const store = createStore(
//     crudCreatorApp,
//     applyMiddleware(thunkMiddleware)
// );

// store.subscribe(() => console.log("STATE: ", store.getState()));

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
