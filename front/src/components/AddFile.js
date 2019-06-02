import React, {Component} from 'react';
import {Alert, Button, Col, Form, FormGroup, FormLabel, Row} from "react-bootstrap";
import jwtDecode from "jwt-decode";
import {parse_json, updater} from "../tools";

const ROUTES = [
    { name: "Новинки" },
    { name: "Популярные" },
    { name: "Из фильмов и игр" },
    { name: "Местные исполнители" },
    { name: "Школы игры" },
];

class AddFile extends Component {
    constructor(props) {
        super(props);
        this.url = "http://127.0.0.1:9090/file";
        const token = localStorage.getItem("accessToken");
        if (token === null || token === "") {
            this.state = {
                validated: false,
                file: "",
                category: ROUTES[0].name,
                musician: "",
                filename: "",
                name: "",
                authorized: false,
            };
            return
        }
        let tokenData = jwtDecode(token);
        let interval = (tokenData.exp - (Date.now().valueOf() / 1000))-10;
        if (interval < 0) {
            localStorage.setItem("accessToken", "");
            localStorage.setItem("refreshToken", "");
            localStorage.setItem("login", "");
            this.state = {
                validated: false,
                file: "",
                category: ROUTES[0].name,
                musician: "",
                filename: "",
                name: "",
                authorized: false,
            }
        } else {
            this.state = {
                validated: false,
                file: "",
                category: ROUTES[0].name,
                musician: "",
                filename: "",
                name: "",
                authorized: true,
            };
        }
        //
        // this.state = {
        //     validated: false,
        //     file: "",
        //     category: ROUTES[0].name,
        //     musician: "",
        //     filename: "",
        //     name: "",
        // };
    }

    handleSubmit(event) {
        const form = event.currentTarget;
        if (form.checkValidity() === false) {
            event.preventDefault();
            event.stopPropagation();
            return;
        }
        if (this.state.file === "") {
            alert("Вы не выбрали файл");
            event.preventDefault();
            event.stopPropagation();
            return;
        }
        this.setState({ validated: true });
        const data = JSON.stringify({
            filename: this.state.filename,
            song: this.state.name,
            musician: this.state.musician,
            category: this.state.category,
            content: this.state.file,
        });
        const accessToken = localStorage.getItem("accessToken");
        fetch(this.url + "?access_token="+accessToken, {
            method: "post",
            mode: 'no-cors',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: data
        })
            .then( res => {
                debugger;
                if (res.status === 200) {
                    return parse_json(res);
                } else {
                    return res.json();
                }
            })
            .then(json => {
                debugger;
                if (json.error) {
                    throw new Error(json.error);
                }
                alert("OK!");
            })
            .catch((error) => {
                alert("Проблемы с доступом в джойказино: " + error.message);
            });
    }

    handleFileSelect = evt => {
        let file = evt.target.files[0]; // FileList object
        let callBack = (str) => this.setState({file: str});
        if (typeof file === 'undefined') {
            document.getElementById('fileUploadName').innerHTML = "Выберете файл";
            callBack("");
            this.setState({filename: ""});
            return;
        }
        this.setState({filename: file.name});
        let reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = function () {
            let a = reader.result;
            a = a.split(',')[1];
            callBack(a);
        };
        reader.onerror = function (error) {
            console.log('Error: ', error);
        };
        document.getElementById('fileUploadName').innerHTML = file.name;
    };

    handleCategorySelect = () => {
        const a = document.getElementById('selectCategory').value;
        this.setState({ category: a });
    };

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    render() {
        let body;

        const { validated } = this.state;
        if (this.state.authorized) {
            body = <Form
                noValidate
                validated={validated}
                onSubmit={e => this.handleSubmit(e)}
            >
                <FormLabel>
                    <h3>Загрузка файла</h3>
                </FormLabel>
                <Row>
                    <Col>
                        <FormGroup controlId="musician" >
                            <Form.Label>ИСПОЛНИТЕЛЬ</Form.Label>
                            <Form.Control
                                required
                                value={this.state.musician}
                                onChange={this.handleChange}
                            />`
                        </FormGroup>
                    </Col>
                    <Col>
                        <Form.Label>ФАЙЛ</Form.Label>
                        <div className="input-group">
                            <div className="input-group-prepend">
                                <span className="input-group-text" id="inputGroupFileAddon01">
                                  Загрузка
                                </span>
                            </div>
                            <div className="custom-file">
                                <input
                                    type="file"
                                    className="custom-file-input"
                                    id="inputGroupFile01"
                                    aria-describedby="inputGroupFileAddon01"
                                    onChange={this.handleFileSelect}
                                />
                                <label id="fileUploadName" className="custom-file-label" htmlFor="inputGroupFile01">
                                    Выберете файл
                                </label>
                            </div>
                        </div>
                    </Col>
                </Row>
                <Row>
                    <Col>
                        <FormGroup controlId="name" >
                            <Form.Label>НАЗВАНИЕ</Form.Label>
                            <Form.Control
                                required
                                value={this.state.name}
                                onChange={this.handleChange}
                            />
                        </FormGroup>
                    </Col>
                    <Col>
                        <Form.Label>КАТЕГОРИЯ</Form.Label>
                        <select
                            className="browser-default custom-select"
                            id="selectCategory"
                            onChange={this.handleCategorySelect}
                        >
                            {ROUTES.map((rout, index) => (
                                <option
                                    value={rout.name}
                                >
                                    {rout.name}
                                </option>
                            ))}
                        </select>
                        {/*<Form.Control required/>*/}
                    </Col>
                </Row>
                {/*<Row>*/}
                {/*    <Col>*/}
                {/*        <Form.Label>ИНСТРУМЕНТЫ</Form.Label>*/}
                {/*        <Form.Control/>*/}
                {/*    </Col>*/}
                {/*</Row>*/}
                <Button
                    block
                    bsSize="large"
                    type="submit"
                >
                    Добавить
                </Button>
            </Form>
        } else {
            body = <Alert variant="danger">
                <Alert.Heading>Вы не авторизованы!</Alert.Heading>
                <p>
                    Для добавления файлов на сайт необходимо авторизоваться.
                </p>
            </Alert>
        }

        return(
            <div>
            {body}
            </div>
        );
    }
}

export default AddFile;