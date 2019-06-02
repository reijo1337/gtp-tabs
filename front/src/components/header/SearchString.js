import React, {Component} from 'react';
import {Button, FormControl, InputGroup, Navbar} from "react-bootstrap";

class SearchString extends Component{
    render() {
        return (
            <Navbar className="justify-content-center" bg="dark">
                <InputGroup className="mb-3">
                    <InputGroup.Prepend>
                        <select className="browser-default custom-select">
                            <option value="1">По табулатурам</option>
                            <option value="2">По авторам</option>
                        </select>
                    </InputGroup.Prepend>
                    <FormControl
                        placeholder="Поиск табулатур"
                        aria-label="Recipient's username"
                        aria-describedby="basic-addon2"
                    />
                    <InputGroup.Append>
                        <Button variant="outline-secondary">Найти</Button>
                    </InputGroup.Append>
                    <Button variant="outline-success" href="/upload">Добавить табулатуру</Button>
                </InputGroup>
            </Navbar>
        );
    }
}

export default SearchString;