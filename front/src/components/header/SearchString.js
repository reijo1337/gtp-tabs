import React, {Component} from 'react';
import {Button, FormControl, InputGroup, Navbar} from "react-bootstrap";

class SearchString extends Component{
    constructor(props) {
        super(props);
        this.state = {
            searchType: 1,
            query: "",
        }
    }
    render() {
        return (
            <Navbar className="justify-content-center" bg="dark">
                <InputGroup className="mb-3" >
                    <InputGroup.Prepend>
                        <select className="browser-default custom-select" id="typeID">
                            <option value="1">По табулатурам</option>
                            <option value="2">По авторам</option>
                        </select>
                    </InputGroup.Prepend>
                    <FormControl
                        placeholder="Поиск табулатур"
                        aria-label="Recipient's username"
                        aria-describedby="basic-addon2"
                        id="query"
                    />
                    <InputGroup.Append>
                        <Button variant="outline-secondary" onClick={this.search}>Найти</Button>
                    </InputGroup.Append>
                    <Button variant="outline-success" href="/upload">Добавить табулатуру</Button>
                </InputGroup>
            </Navbar>
        );
    }

    search = () => {
        const e = document.getElementById("typeID");
        const selectType = e.options[e.selectedIndex].value;
        const qe = document.getElementById("query").value;
        if (selectType === "1") {
            window.location = "http://127.0.0.1:3000/tabs/"+qe;
        } else {
            window.location = "http://127.0.0.1:3000/musicians/"+qe;
        }
    }
}

export default SearchString;