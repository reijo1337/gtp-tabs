import React, {Component} from 'react';
import {Button} from "react-bootstrap";

class TabWithSize extends Component{
    constructor(props) {
        super(props);
        let {data} = this.props;
        this.data = data;
    }

    render() {
        return (
            <div className="container">
                <Button variant="primary" size="lg" block href={"/post/"+this.data.id}>
                    {this.data.musician} - {this.data.name}; Размер: {this.data.size}
                </Button>
            </div>
        );
    }
}

export default TabWithSize;