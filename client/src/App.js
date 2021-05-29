import React from 'react';
import './App.css';
import Home from './Home';
import Listing from './Listing';
import {Route, Switch } from 'react-router-dom';

export default function App() {

  return (
    <main>
      <Switch>
        <Route exact path="/" component={Home} />
        <Route exact path="/discover" component={Listing} />
      </Switch>
    </main>
  );
}
