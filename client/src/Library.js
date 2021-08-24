import React from 'react';
import ComingSoon from './ComingSoon.js';
import Layout from './Layout.js';
import './Layout.css';

export default class Library extends React.Component {
    render(){
        return <Layout nav="true">
            <ComingSoon />
        </Layout>
    };
}