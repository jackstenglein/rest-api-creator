import React from 'react';
import logo from './logo.svg';
import './App.css';
import { createObject } from './actions/actions.js';
import fetch from 'cross-fetch';
import { BrowserRouter as Router, Route, Link } from 'react-router-dom';

function App(props) {

    fetch('http://localhost:8000/app/login', {
        method: "PUT",
        body: JSON.stringify({"username": "jack3", "password": "jackpassword"}),
        headers: {
            "Content-Type": "application/json"
        },
        credentials: "include"
    }).then(function(response, error) {
        console.log("Login Response: ", response);
        console.log("Headers: ", response.headers);

        if (error) {
            console.log('An error ocurred.', error);
        }
        return response.json();
    }).then(function(json) {
        console.log(props.store.getState());
        const unsubscribe = props.store.subscribe(() => console.log("STATE: ", props.store.getState()));

        props.store.dispatch(createObject({"name": "Test 7", "attributes": []}));
    });

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
