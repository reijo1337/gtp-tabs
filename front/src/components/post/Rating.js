import React, {Component} from 'react';
import StarRatingComponent from 'react-star-rating-component';
import {parse_json} from "../../tools";


class Rating extends Component {
    constructor(props) {
        super(props);
        this.url = "http://localhost:9090/rating";
        let {rating} = this.props;
        this.post_id = this.props.post_id;
        this.state = {
            rating: rating,
        };
    }

    sendRating = async (data) => {
        const accessToken = localStorage.getItem("accessToken");
        let response = await fetch(this.url + "?access_token=" + accessToken, {
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
                alert("ok");
                return json;
            })
            .catch(error => {
                alert("Проблемы с доступом в джойказино: " + error.message);
                return error;
            });
        return response;
    };

    onStarClick(nextValue, prevValue, name) {
        const data = JSON.stringify({
            post_id: this.post_id,
            rating: nextValue,
        });
        this.sendRating(data);
        this.setState({rating: nextValue});
    }

    render() {
        const { rating } = this.state;

        return (
            <div>
                <StarRatingComponent
                    name="rate1"
                    starCount={5}
                    value={rating}
                    onStarClick={this.onStarClick.bind(this)}
                />
            </div>
        );
    }
}

export default Rating;