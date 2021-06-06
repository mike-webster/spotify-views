import React from 'react';
import Home from './Home';
import ActionSelect from './ActionSelect';
import WordCloud from './WordCloud';
import Library from './Library';
import {Route, Switch } from 'react-router-dom';

export default function Router() {
  return (
    <main>
      <Switch>
        <Route exact path="/" component={Home} />
        <Route exact path="/discover" component={ActionSelect} />
        <Route exact path="/wordcloud" component={WordCloud} />
        <Route exact path="/library" component={Library} />
      </Switch>
    </main>
  );
}
