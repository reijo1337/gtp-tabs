import React, {Component} from 'react';
import {parse_json} from "../../tools";
import {Alert, ListGroup, ListGroupItem} from "react-bootstrap";
import Rating from "./Rating";
import Comment from "./Comment";

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
            const arrearsList = this.comments.map(ar =>
                <ListGroupItem key={ar.id}>
                    <Comment comment={ar}/>
                </ListGroupItem>
            );
            body = <div>
                <h1>{this.musician_name} - {this.song_name}</h1>
                <p>
                    Рейтинг: <Rating rating={this.rating} post_id={this.post_id}/>
                </p>
                <p>
                    Размер: {this.size} байт
                </p>
                <p>
                    Скачать: <a href={this.download} target="_blank">{this.filename}</a>
                </p>
                <p>Комментарии:</p>
                <ListGroup>
                    {arrearsList}
                </ListGroup>
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
                });
        }
    }
}

export default Post;