import React, {Component} from 'react';
import {Button, Form} from "react-bootstrap";
import {parse_json} from "../../tools";

class AddComment extends Component {
    constructor(props) {
        super(props);
        let {postID} = this.props;
        this.postID = postID;
        this.state = {
            text: "",
        }
    }

    validateForm() {
        return this.state.text.length > 0;
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleSubmit = () => {
        const profileID = localStorage.getItem("profileID");
        const accessToken = localStorage.getItem("accessToken");
        const url = "http://127.0.0.1:9090/post/" + this.postID + "?access_token=" + accessToken;
        let data = JSON.stringify({
            author_id: parseInt(profileID),
            content: this.state.text,
        });
        this.sendData(data, url);
    };

    sendData = async (data, url) => {
        let response = await fetch(url, {
            method: "post",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: data
        })
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
                window.location.href = "http://127.0.0.1:3000/post/" + json.tab.id;
                return json;
            })
            .catch(error => {
                alert("Проблемы с доступом в джойказино: " + error.message);
                return error;
            });
        return response;
    };

    render() {
        return (
            <div>
                <Form
                    onSubmit={this.handleSubmit}
                >
                    <Form.Group controlId="text">
                        <Form.Label>Добавить комментарий</Form.Label>
                        <Form.Control
                            as="textarea" rows="3"

                            value={this.state.text}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Button
                        variant="primary"
                        disabled={!this.validateForm()}
                        type="submit"
                    >Добавить комментарий</Button>
                </Form>
            </div>
        );
    }
}

export default AddComment;