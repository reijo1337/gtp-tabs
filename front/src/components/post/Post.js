import React, {Component} from 'react';
import {parse_json} from "../../tools";
import {Alert} from "react-bootstrap";
import Rating from "./Rating";
import {Link} from "react-router-dom";

class Post extends Component{
    constructor(props){
        super(props);
        this.profile_id = this.props.match.params.id;
        this.url = "http://localhost:9090/post/" + this.profile_id;
        this.state = {
            isLoaded: false,
        };
        this.loadPost();
    }

    render() {
        let body;
        if (!this.state.isLoaded) {
            body = <Alert variant="danger">
                <Alert.Heading>Ошибка!</Alert.Heading>
                <p>
                    Произошла ошибка. Попробуйте обновить страницу.
                </p>
            </Alert>
        } else {
            body = <div>
                <h1>{this.musician_name} - {this.song_name}</h1>
                <p>
                    Рейтинг: <Rating rating={this.rating} post_id={this.post_id}/>
                </p>
                <p>
                    Размер: {this.size}
                </p>
                <p>
                    Скачать:
                    <a href={this.download} target="_blank">{this.filename}</a>

                    {/*{this.download}*/}
                </p>
            </div>
        }
        return (
            <div className="d-block mx-auto">
                {body}
            </div>
        );
    }

    loadPost = () => {
        if (!this.state.isLoaded) {
            fetch(this.url)
                .then(res => {
                    debugger;
                    if (res.status === 200) {
                        return parse_json(res);
                    } else {
                        return res.json();
                    }
                })
                .then(json => {
                    if (json.error) {
                        throw new Error(json.error);
                    }
                    debugger;
                    this.comments = json.post.comments;

                    this.song_name = json.post.song_name;
                    this.rating = json.post.rating;
                    this.author_id = json.post.author_id;
                    this.post_id = json.post.id;
                    this.musician_id = json.tab.musician.id;
                    this.musician_name = json.tab.musician.name;
                    this.size = json.tab.size;
                    this.download = "http://localhost:9090/file?name="+json.tab.name;
                    this.filename = json.tab.name;

                    this.setState({
                        isLoaded: true,
                    });

                })
                .catch((error) => {
                    alert("Cant get arrears: " + error.message);
                });
        }
    }
}

export default Post;