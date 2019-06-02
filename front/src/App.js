import React, {Component} from 'react';
import Header from "./components/header/Header";
import Glagna from "./components/glagna/Glagna";
import { Route } from 'react-router-dom';
import AddFile from "./components/AddFile";
import Post from "./components/post/Post";

class App extends Component {
  render(){
  return (
    <div>
      <Header/>
      <div className="container">
      <Route exact path="/" component={Glagna}/>
      <Route path="/upload" component={AddFile}/>
      <Route path="/post/:id" component={Post}/>
      </div>
    </div>
  );
  }

}

export default App;
