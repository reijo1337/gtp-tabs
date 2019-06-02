import React, {Component} from 'react';
import Header from "./components/header/Header";
import Glagna from "./components/glagna/Glagna";
import { Route } from 'react-router-dom';
import AddFile from "./components/AddFile";
import Post from "./components/post/Post";
import CategorySearch from "./components/search/CategorySearch";
import MusicianSearch from "./components/search/MusicianSearch";

class App extends Component {
  render(){
  return (
    <div>
      <Header/>
      <div className="container">
          <Route exact path="/" component={Glagna}/>
          <Route path="/upload" component={AddFile}/>
          <Route path="/post/:id" component={Post}/>
          <Route path="/category/:name" component={CategorySearch}/>
          <Route path="/musician/:id" component={MusicianSearch}/>
      </div>
    </div>
  );
  }

}

export default App;
