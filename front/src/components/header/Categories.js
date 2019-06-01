import React, {Component} from 'react';
import {Nav} from "react-bootstrap";

const ROUTES = [
    { name: "Новинки" },
    { name: "Популярные" },
    { name: "Из фильмов и игр" },
    { name: "Местные исполнители" },
    { name: "Школы игры" },
];

class Categories extends Component {
    render() {
        return (
            <div>
                <Nav className="justify-content-center">
                    {ROUTES.map((rout, index) => (
                        <Nav.Item key={rout.name}>
                            <Nav.Link key={rout.name} eventKey={index}>{rout.name}</Nav.Link>
                        </Nav.Item>
                    ))}
                </Nav>
            </div>
        );
    }
}

export default Categories;