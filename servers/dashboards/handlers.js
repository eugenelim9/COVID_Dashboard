const allDashHandler = async (req, res, { Dashboard }) => {
    try {
        const allDashboards = await Dashboard.find({"private":false})
        res.set('Content-Type', 'application/json')
        res.json(allDashboards)

    } catch(err) {
        res.status(500).send("internal server error")
        return
    }
}

const specificDashGetHandler = async (req, res, { Dashboard }, user) => {
    const dashboardID = req.params.dashID

    if (dashboardID.length !== 24 && dashboardID !== "me") {
        res.status(400).send('dashboard not found');
        return;
    }

    try {
        // try to find dashboard
        let theDash
        if (dashboardID === "me") {
            theDash = await Dashboard.find({"creator.id":user.id})
        } else {
            theDash = await Dashboard.findById(dashboardID)
        }
        if (!theDash) {
            res.status(400).send('dashboard not found')
            return
        }

        res.send(theDash)
    } catch(err) {
        res.status(500).send("internal server error")
        return
    }
}

const allDashPostHandler = async (req, res, { Dashboard }, user) => {
    // check that dashboard exists
    try {
        const { title, description, params, private } = req.body
        if (!title) {
            res.status(400).send("must include dashboard title")
            return
        }
        if (!params) {
            res.status(400).send("must include dashboard param")
            return
        }

        const newDashboard = {
            title: title,
            creator: user,
            description: description,
            params: params,
            private: private,
            createdAt: new Date()
        }
        
        
        const exists = await Dashboard.exists({"creator.id":newDashboard.creator.id})
        if (exists) {
            res.status(400).send("creator already has dashboard")
            return
        }

        const query = new Dashboard(newDashboard)
        const result = await query.save() 

        res.set('Content-Type', 'application/json')
        res.status(201).send(result)

    } catch(err) {
        res.status(500).send("internal server error")
        return
    }
}

const specificDashPatchHandler = async (req, res, { Dashboard }, user) => {
    const dashboardID = req.params.dashID
    if (dashboardID !== "me") {
        res.status(400).send('dashboard not found');
        return;
    }
    try {
        const theDash = await Dashboard.find({"creator.id":user.id})
        if (theDash.length === 0) {
            res.status(400).send("dashboard not found")
            return
        }
        // check authorization
        // if (theDash.creator.id !== user.id) {
        //     res.status(403).send("user not authorized")
        //     return
        // }
        let {title, description, params, private} = req.body
        if (!title) {
            title = theDash[0].title
        }

        if (!description) {
            description = ""
        }

        if (private === undefined) {
            private = theDash[0].private
        }

        if (!params || JSON.stringify(params) === "{}") {
            params = theDash[0].params
        }
        const updatedDash = await Dashboard.findByIdAndUpdate(theDash[0]._id,
            {
                description: description, 
                editedAt: new Date(), 
                private:private,
                title: title,
                params: params
            },
            {new: true}
        )
        res.set('Content-Type', 'application/json')
        res.status(201).json(updatedDash)

    } catch(err) {
        res.status(500).send("internal server error")
        return
    }
}

const dataGetHandler = async (req, res, { CountriesCovid }) => {
    try {
        const query = CountriesCovid.where({})
        data = await query.find()
        res.send(data)
    } catch(e) {
        res.status(500).send('internal server error')
        return
    }
}

const dataDelHandler = async (req, res, { CountriesCovid }) => {
    try {
        const query = CountriesCovid.where({})
        data = await query.remove()
        res.send("successfully removed")
    } catch(e) {
        res.status(500).send('internal server error')
        return
    }
    
}

module.exports = { 
    allDashHandler, 
    specificDashGetHandler,
    allDashPostHandler,
    specificDashPatchHandler,
    dataGetHandler,
    dataDelHandler
}