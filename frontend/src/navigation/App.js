import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import PrimaryLayout from './PrimaryLayout.js';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Switch>
        {/* <Route path="/auth" component="" /> */}
        <Route path="/app" component={PrimaryLayout} />
      </Switch>
    </BrowserRouter>
  );
}
 
export default App;
