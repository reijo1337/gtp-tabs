import React, {Component} from 'react';
import {parse_json} from "../../tools";
import {ListGroup, ListGroupItem} from "react-bootstrap";
import MusicianWithCount from "./MusicianWithCount";
import TabWithSize from "./TabWithSize";

class CategorySearch extends Component{
    constructor(props) {
        super(props);
        this.name = this.props.match.params.id;
        this.url = "http://localhost:9090/musician/"+this.name;
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
                    <TabWithSize data={ar}/>
                </ListGroupItem>
            );
            body = <div>
                    <h1>{this.data.musician}</h1>
                    <ListGroup>
                        {resList}
                    </ListGroup>
                </div>
        }
        return (
            <div>
                {body}
            </div>
        );
    }
}

export default CategorySearch;