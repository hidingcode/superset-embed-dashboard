## Superset Dashboard Embedded Bridge

A full-stack implementation for embedding Apache Superset dashboards into a custom web application. This project features a secure Backend Auth API that handles Guest Token exchange, ensuring your Superset credentials remain private while providing a seamless user experience.

## Tech Stack
- Frontend: React ([Devias Kit](https://github.com/devias-io/material-kit-react 'Devias Kit - React')) + `@superset-ui/embedded-sdk`
- Backend: Python (Flask) / Go
- Visualizations: Apache Superset

## Architecture Diagram

<img width="1075" height="373" alt="image" src="https://github.com/user-attachments/assets/7eac47fa-a2ca-43a8-8840-35c4d68d0d59" />

## File Structure
```
┌── backend
    ├── Flask
        ├──.venv
        ├──app.py
    ├── Go
        ├──go.mod
        ├──go.sum
        ├──main.go
├── frontend
	├── .editorconfig
  ├── .eslintrc.js
  ├── .gitignore
  ├── CHANGELOG.md
  ├── LICENSE.md
  ├── next-env.d.ts
  ├── next.config.js
  ├── package.json
  ├── README.md
  ├── tsconfig.json
  ├── public
  └── src
  	├── components
  	├── contexts
  	├── hooks
  	├── lib
  	├── styles
  	├── types
  	└── app
  		├── layout.tsx
  		├── page.tsx
  		├── auth
  		└── dashboard
└── README.md
```

