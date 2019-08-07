import React from 'react';
import Link from 'react-router-dom/Link';

const MENU_ITEMS = [
    {
        id: 'projectDetails',
        name: 'Project Details',
        link: '/details'
    },
    {
        id: 'objects',
        name: 'Objects',
        link: '/objects'
    },
    {
        id: 'endpoints',
        name: 'Endpoints',
        link: '/endpoints'
    },
    {
        id: 'docs',
        name: 'Documentation',
        link: '/documentation'
    },
    {
        id: 'tests',
        name: 'Tests',
        link: '/tests'
    }
];

const SELECTED_CLASS = 'project-menu-link-selected';
const DESELECTED_CLASS = 'project-menu-link-deselected';

function ProjectMenu(props) {
    return (
        <div>
            {
                MENU_ITEMS.map(function(item) {
                    let className = item.id === props.selectedId ? SELECTED_CLASS : DESELECTED_CLASS;
                    return (
                        <Link key={item.id} to={item.link} className={className}>
                            <div className='project-menu-item'>{item.name}</div>
                        </Link>
                    );
                })
            }
        </div>
    );
}

export default ProjectMenu;
