import { Grid, Paper, Typography } from "@mui/material";
import { DndSender, ItemTypes } from "./dnd";

const MacroPanel = ({ macros }) => {
    return <div style={{
        height: "261px",
        position: "relative",
        padding: "24px",
        display: "inline-flex",
    }}>
        <Grid
            container
            direction="column"
            justifyContent="flex-start"
            alignItems="flex-start"
            spacing={"5px"}
            sx={{
                width: "auto",
                height: "100%",
                // marginLeft: "70px"
            }}
        >
            {Object.keys(macros).map((value) => (
                <Grid key={value} item xs={2} >
                    <DndSender args={{ ...macros[value], key: value }} type={ItemTypes.CARD} key={value}>
                        <Paper
                            elevation={3}
                            sx={{
                                width: "auto",
                                height: "100%",
                                backgroundColor: "button.mouse.main",
                                "&:hover": {
                                    backgroundColor: "button.mouse.hover",
                                },
                                color: "button.mouse.text",
                                padding:"4px"
                            }} >
                            <Typography variant="h6" component="h1" noWrap={true}>
                                {macros[value]["name"]}
                            </Typography>
                            <Typography variant="body2" component="h1" noWrap={true}>
                                {macros[value]["description"]}
                            </Typography>
                        </Paper>
                    </DndSender>
                </Grid>
            ))}
            

        </Grid>
    </div>
}


export default MacroPanel;