import React, {Component} from 'react';
import TopHeader from "./TopHeader";
import Categories from "./Categories";

class Header extends Component {
    render() {
        return (
            <div>
                <TopHeader/>
                <Categories/>
            </div>
        );
    }
}

export default Header;