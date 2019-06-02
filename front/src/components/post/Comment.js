import React, {Component} from 'react';
import {parse_json} from "../../tools";
import {Badge} from "react-bootstrap";

class Comment extends Component {
    constructor(props) {
        super(props);
        const {comment} = this.props;
        this.comment = comment;
        this.url = "http://127.0.0.1:9090/profile/"+this.comment.author_id;
        this.state = {
            isLoaded: false,
        };
        this.loadInfo();
    }

    loadInfo = () => {
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
                    this.author_name = json.name;

                    this.setState({
                        isLoaded: true,
                    });

                })
                .catch((error) => {
                });
        }
    };

    render() {
        let body = "";
        if (this.state.isLoaded) {
            body = <div>
                <Badge pill variant="dark">
                    {this.author_name}
                </Badge>
                <p>
                    {this.comment.content}
                </p>
            </div>
        }
        return(
            <div className="container">
                {body}
            </div>
        );
    }
}

export default Comment;