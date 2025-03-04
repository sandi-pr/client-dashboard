<script>
    let configModal = [
    {
      header: {
        iconLeft: {
          icon: "home",
          color: "primary",
        },
        contentCenter: [
          {
            text: {
              text: "Add Document",
              color: "primary",
              preset: "Text md/Semibold",
            },
          },
          {
            text: {
              text: "Home",
              color: "primary",
              preset: "Text md/Semibold",
            },
          },
        ],
        iconRight: {
          icon: "home",
          color: "primary",
        },
      },
    },
    {
      divider: {
        className: "",
      }
    },
    {
      mainContent: [
        {
          tabBar: {
            tabDatas: {
              selected: "information",
              data: [
                {
                  name: "information",
                },
                {
                  name: "recipients",
                }
              ],
              changeTabbar: () => {}
            }
          }
        },
        {
          text: {
            text: "Home",
            color: "primary",
            preset: "Text md/Semibold",
          },
        },
        {
          input: {
            type: "text",
            onEdit: false,
            value: "",
            oninput: () => {}
          }
        },
        {
          text: {
            text: "Label",
            color: "primary",
            preset: "Text md/Semibold",
          }
        },
        {
          uploadFile: {
            onUpload: () => {},
            onDrop: () => {},
            files: [],
          }
        },
        {
          text: {
            text: "Label",
            color: "primary",
            preset: "Text md/Semibold",
          },
        },
        {
          multiSelect: {
            alltags: [],
            tags: [],
          }
        },
      ]
    },
    {
      footer: [
        {
          button: {
            text: "Add Document",
            onClick: () => {}
          },
        }
      ]
    }
  ]
</script>