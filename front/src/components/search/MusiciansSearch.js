import React, {Component} from 'react';
import {parse_json} from "../../tools";
import MusicianWithCount from "./MusicianWithCount";
import {ListGroup, ListGroupItem} from "react-bootstrap";

class MusiciansSearch extends Component{
    constructor(props) {
        super(props);
        this.name = this.props.match.params.name;
        this.url = "http://localhost:9090/musicians/"+this.name;
        this.state = {
            isLoaded: false,
        };
        this.loadCategory();
    }

    loadCategory = () => {
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
                    this.data = json;

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
            const resList = this.data.map(ar =>
                <ListGroupItem key={ar.id}>
                    <MusicianWithCount data={ar}/>
                </ListGroupItem>
            );
            body = <ListGroup>
                {resList}
            </ListGroup>
        }
        return (
            <div>
                <h1>Результаты по {this.name}</h1>
                {body}
            </div>
        );
    }
}

export default MusiciansSearch;