import React, {Component} from 'react';
import {Button} from "react-bootstrap";

class MusicianWithCount extends Component{
    constructor(props) {
        super(props);
        let {data} = this.props;
        this.data = data;
    }

    render() {
        return (
            <div className="container">
                <Button variant="primary" size="lg" block href={"/musician/"+this.data.id}>
                    {this.data.name}; Файлов: {this.data.count}
                </Button>
            </div>
        );
    }
}

export default MusicianWithCount;