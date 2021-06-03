import React from 'react';
import './App.css';
import Home from './Home';
import ActionSelect from './ActionSelect';
import WordCloud from './WordCloud';
import {Route, Switch } from 'react-router-dom';

export default function App() {

  return (
    <main>
      <Switch>
        <Route exact path="/" component={Home} />
        <Route exact path="/discover" component={ActionSelect} />
        <Route exact path="/wordcloud" component={WordCloud} />
      </Switch>
    </main>
  );
}
